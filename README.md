# Mesh-test

Sandbox for testing Open Service Mesh configuration.

After installing [OSM](https://openservicemesh.io/) I want to be able to prove out some of the basic use cases I am looking to leverage.

Use cases:
* service-to-service encryption via envoy sidecars.
* access policy (limit what pods can and can't do within the k8s network).
* monitoring.

For this purpose I am going to create a small front end web service - mesh-server that will consume information from a back-end microservice mesh-service. To keep things lightweight I'm going to write and deploy these services using GO.

## Service overview

mesh-service listens on port 8081 and provides a random number to each request

mesh server listens on port JB_MESH_SERVER_PORT (or a default of :8080), and polls JB_MESH_SERVICE (default http://localhost:8080) for each request and dispays a formatted result. 
