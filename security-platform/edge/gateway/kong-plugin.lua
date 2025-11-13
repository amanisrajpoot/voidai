-- Security Platform Kong Plugin
-- Provides integration with Security Platform control plane

local BasePlugin = require "kong.plugins.base_plugin"
local http = require "resty.http"

local SecurityPlatformHandler = BasePlugin:extend()

SecurityPlatformHandler.PRIORITY = 1000
SecurityPlatformHandler.VERSION = "1.0.0"

function SecurityPlatformHandler:new()
  SecurityPlatformHandler.super.new(self, "security-platform")
end

function SecurityPlatformHandler:access(conf)
  SecurityPlatformHandler.super.access(self)
  
  local control_plane_url = conf.control_plane_url or "https://api.securityplatform.com"
  local auth_key = conf.auth_key
  
  -- Get request metadata
  local request_id = ngx.var.request_id or ngx.var.remote_addr .. "-" .. ngx.now()
  ngx.req.set_header("X-Request-ID", request_id)
  
  local request_metadata = {
    request_id = request_id,
    method = ngx.req.get_method(),
    uri = ngx.var.request_uri,
    headers = ngx.req.get_headers(),
    remote_addr = ngx.var.remote_addr,
    user_agent = ngx.var.http_user_agent,
    timestamp = ngx.now(),
  }
  
  -- Check policy with control plane
  local httpc = http.new()
  httpc:set_timeout(100) -- 100ms timeout
  
  local res, err = httpc:request_uri(control_plane_url .. "/api/v1/policy/check", {
    method = "POST",
    headers = {
      ["Content-Type"] = "application/json",
      ["Authorization"] = "Bearer " .. auth_key,
    },
    body = ngx.json.encode({
      action = "http_request",
      context = request_metadata,
    }),
  })
  
  if res and res.status == 200 then
    local policy_result = ngx.json.decode(res.body)
    if not policy_result.allowed then
      ngx.status = 403
      ngx.say(ngx.json.encode({
        error = "Request blocked by security policy",
        reason = policy_result.reason,
      }))
      ngx.exit(403)
    end
  end
  
  -- Send event to control plane (async, non-blocking)
  ngx.timer.at(0, function()
    local httpc_async = http.new()
    httpc_async:set_timeout(1000)
    
    httpc_async:request_uri(control_plane_url .. "/api/v1/events", {
      method = "POST",
      headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer " .. auth_key,
      },
      body = ngx.json.encode({
        type = "http_request",
        data = request_metadata,
      }),
    })
  end)
end

function SecurityPlatformHandler:response(conf)
  SecurityPlatformHandler.super.response(self)
  
  -- Capture response metadata
  local response_metadata = {
    status_code = ngx.status,
    headers = ngx.resp.get_headers(),
    response_time = ngx.var.request_time,
  }
  
  -- Send to control plane (async)
  ngx.timer.at(0, function()
    local control_plane_url = conf.control_plane_url or "https://api.securityplatform.com"
    local auth_key = conf.auth_key
    
    local httpc = http.new()
    httpc:set_timeout(1000)
    
    httpc:request_uri(control_plane_url .. "/api/v1/events", {
      method = "POST",
      headers = {
        ["Content-Type"] = "application/json",
        ["Authorization"] = "Bearer " .. auth_key,
      },
      body = ngx.json.encode({
        type = "http_response",
        data = response_metadata,
      }),
    })
  end)
end

return SecurityPlatformHandler
