# Azure CLI Proxy for Docker

This is a docker image that allows you to run the Azure CLI in a container as a proxy.

With this you could run the Azure CLI in a container and use it to manage your Azure resources without having to install the Azure CLI on your local machine.

It also mimics the idea of a User Managed Identity by using the Azure CLI to authenticate with Azure and then using the Azure CLI to get a token that can be used to authenticate with Azure services.
