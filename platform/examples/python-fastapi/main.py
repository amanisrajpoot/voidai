"""
FastAPI example with Security Platform agent
"""

from fastapi import FastAPI
from security_platform import init_agent, SecurityMiddleware, AgentConfig

# Initialize agent
config = AgentConfig(
    control_plane_url="http://localhost:3000",
    auth_key="demo-key",
    service_name="fastapi-demo",
    environment="development",
    redaction_rules=["password", "token", "secret"],
)

agent = init_agent(config)

# Create FastAPI app
app = FastAPI(title="Security Platform Demo API")

# Add security middleware
app.add_middleware(SecurityMiddleware)


@app.get("/")
def read_root():
    return {"message": "Hello World", "service": "fastapi-demo"}


@app.get("/users/{user_id}")
def get_user(user_id: int):
    # This request will be automatically captured by the middleware
    return {"user_id": user_id, "name": "John Doe"}


@app.post("/login")
def login(username: str, password: str):
    # Password will be automatically redacted
    agent.capture_event("login_attempt", {
        "username": username,
        "password": password,  # Will be redacted
    })
    return {"message": "Login successful", "username": username}


@app.get("/admin")
def admin_panel():
    # Suspicious path - will trigger security event
    return {"message": "Admin panel"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
