"""
Security Platform Python Agent - FastAPI Quick Start
"""

import os
from fastapi import FastAPI
from security_platform_agent import Agent, create_middleware

# Initialize agent
agent = Agent({
    "control_plane_url": os.getenv(
        "SECURITY_PLATFORM_CONTROL_PLANE_URL", "https://api.securityplatform.com"
    ),
    "auth_key": os.getenv("SECURITY_PLATFORM_AUTH_KEY"),
    "service_name": "my-api-service",
    "environment": os.getenv("ENVIRONMENT", "production"),
    "otlp_endpoint": os.getenv("SECURITY_PLATFORM_OTLP_ENDPOINT"),
    "policy": {
        "mode": "observe",  # Start in observe mode
        "auto_enable_blocking": False,
    },
    "redaction_rules": [
        {"pattern": "password|passwd", "field": "*"},
        {"pattern": r"\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}", "field": "*"},  # Credit card
    ],
})

# Start agent
agent.start()

app = FastAPI()

# Add security platform middleware
app.add_middleware(
    create_middleware(agent, mode="observe", capture_response_body=False)
)


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/api/login")
async def login(request: Request):
    # Request body will be automatically redacted
    data = await request.json()
    username = data.get("username")
    password = data.get("password")  # Will be redacted in telemetry

    # Check policy
    policy_result = agent.check_policy(
        "login",
        {"username": username, "ip": request.client.host},
    )

    if not policy_result["allowed"]:
        raise HTTPException(
            status_code=403, detail="Login blocked by security policy"
        )

    # Your login logic here
    return {"success": True}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
