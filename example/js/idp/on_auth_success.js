// on_auth_success.js
// This script is used to authenticate the user against the federation provider
// ------------
// When using Federated OAuth2/OIDC Login, 
// a script is used to bridge between the user authorization details received from the federated provider
// and the requested authorization details in the originating SMART on FHIR application. 
// This script is used to assign appropriate permissions and inject any other required details into the user session. 
// It may obtain all required information by inspecting the access token details, 
// or it may make additional service calls to fetch information.

println("on_auth_success.js");

function onAuthenticationSuccess(theOutcome, theOutcomeFactory, theContext) {
    
    // This Identity Provider is for Patient Portal. It retrieves the patientId as HDID claim.
    var patientId = theContext.getClaim('hdid');
  
       // Add a log line for troubleshooting
   Log.info("User " + theOutcome.getUsername() + " has authorized for patient: " + patientId + " with scopes: " + theContext.getApprovedScopes());

   // Assign appropriate Smile CDR permissions
   theOutcome.addAuthority('FHIR_READ_ALL_IN_COMPARTMENT', 'Patient/' + patientId);
   theOutcome.addAuthority('FHIR_WRITE_ALL_IN_COMPARTMENT', 'Patient/' + patientId);

   // Now set the launch context to include the patient ID
   theOutcome.addLaunchResourceId('patient', patientId);

   return theOutcome;

}