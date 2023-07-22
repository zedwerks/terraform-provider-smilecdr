// federationAuthScript.js
// This script is used to authenticate the user against the federation provider
// ------------
// When using Federated OAuth2/OIDC Login, 
// a script is used to bridge between the user authorization details received from the federated provider
// and the requested authorization details in the originating SMART on FHIR application. 
// This script is used to assign appropriate permissions and inject any other required details into the user session. 
// It may obtain all required information by inspecting the access token details, 
// or it may make additional service calls to fetch information.

println("federationAuthScript.js");