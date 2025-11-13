.PHONY: install uninstall deploy

# Install CRDs
install:
	kubectl apply -f config/crd/bases/

# Uninstall CRDs
uninstall:
	kubectl delete -f config/crd/bases/

# Deploy operator
deploy:
	kubectl apply -f deploy/operator.yaml

# Build operator image
build:
	docker build -t securityplatform/operator:latest .

# Push operator image
push:
	docker push securityplatform/operator:latest
