from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="security-platform-python-agent",
    version="0.1.0",
    author="Security Platform",
    description="Security Platform Python Agent with RASP capabilities",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.8",
    install_requires=[
        "opentelemetry-api>=1.20.0",
        "opentelemetry-sdk>=1.20.0",
        "opentelemetry-exporter-otlp-proto-http>=1.20.0",
        "opentelemetry-instrumentation-fastapi>=0.42b0",
        "opentelemetry-instrumentation-requests>=0.42b0",
        "fastapi>=0.100.0",
        "starlette>=0.27.0",
    ],
    extras_require={
        "django": ["opentelemetry-instrumentation-django>=0.42b0"],
        "flask": ["opentelemetry-instrumentation-flask>=0.42b0"],
    },
)
