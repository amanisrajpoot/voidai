"""
FastAPI example with Security Platform Python Agent
"""
from fastapi import FastAPI, Request
from security_platform import SecurityMiddleware, AgentConfig, create_agent
import os

app = FastAPI()

# Configure agent
config = AgentConfig(
    control_plane_url=os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL", "http://localhost:8080"),
    auth_key=os.getenv("SECURITY_PLATFORM_AUTH_KEY", ""),
    service_name="fastapi-app",
    environment=os.getenv("ENV", "development"),
    local_policy="observe",
    redaction_rules=["password", "card", "ssn"],
)

# Create agent and add middleware
agent = create_agent(config)
app.add_middleware(SecurityMiddleware, agent=agent)

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.post("/api/users")
async def create_user(request: Request, user_data: dict):
    # Access security context
    context = request.state.security_context
    
    # Capture custom security event
    agent.capture_event(
        "user_created",
        data={"user_id": user_data.get("id")},
        taint_tags=["user_creation"] if user_data.get("admin") else []
    )
    
    return {"status": "created", "request_id": context.request_id}

@app.get("/api/users/{user_id}")
async def get_user(request: Request, user_id: str):
    context = request.state.security_context
    
    # Your business logic here
    return {"user_id": user_id, "trace_id": context.trace_id}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
