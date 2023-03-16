
.PHONY: integrationtests
integrationtests:
	export XMC_MONGODB_URI=mongodb://root:toor123@mongodb:27017
	export XMC_DB_NAME=xmtest
	export XMC_COMPANY_COLLECTION=companies
	export XMC_SIGNING_SECRET=signing_secret
	docker-compose up --force-recreate -d mongodb kafka
	go test -v ./tests  --tags=integration

.PHONY: run
run:
	docker-compose up --force-recreate --build -d