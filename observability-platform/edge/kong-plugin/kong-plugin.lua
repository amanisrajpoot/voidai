local BasePlugin = require "kong.plugins.base_plugin"
local http = require "resty.http"

local ObservabilityPlugin = BasePlugin:extend()

ObservabilityPlugin.PRIORITY = 1000
ObservabilityPlugin.VERSION = "1.0.0"

function ObservabilityPlugin:new()
  ObservabilityPlugin.super.new(self, "observability-plugin")
end

function ObservabilityPlugin:access(conf)
  ObservabilityPlugin.super.access(self)
  
  local otel_endpoint = conf.otel_endpoint or "http://otel-collector:4318"
  
  -- Capture request metadata
  local request_data = {
    method = ngx.var.request_method,
    uri = ngx.var.request_uri,
    headers = ngx.req.get_headers(),
    remote_addr = ngx.var.remote_addr,
    service = ngx.var.upstream_service_name or "unknown",
  }
  
  -- Send to OTEL collector
  local httpc = http.new()
  httpc:set_timeout(1000)
  
  local res, err = httpc:request_uri(otel_endpoint .. "/v1/traces", {
    method = "POST",
    headers = {
      ["Content-Type"] = "application/json",
    },
    body = ngx.json.encode({
      resourceSpans = {{
        resource = {
          attributes = {
            {key = "service.name", value = {stringValue = request_data.service}},
          }
        },
        scopeSpans = {{
          spans = {{
            name = request_data.method .. " " .. request_data.uri,
            kind = "SPAN_KIND_SERVER",
            startTimeUnixNano = ngx.now() * 1000000000,
            attributes = {
              {key = "http.method", value = {stringValue = request_data.method}},
              {key = "http.url", value = {stringValue = request_data.uri}},
            }
          }}
        }}
      }}
    })
  })
  
  if err then
    ngx.log(ngx.ERR, "Failed to send telemetry: ", err)
  end
end

return ObservabilityPlugin
