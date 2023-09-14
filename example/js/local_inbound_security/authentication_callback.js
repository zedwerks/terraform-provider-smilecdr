/**  authentication_callback_script.js */

/**
 * This is an authentication callback script for the
 * SMART Inbound Security module.
 *
 * @param theOutcome The outcome object. This contains details about the user that was created
 * 	in response to the incoming token.
 * 	
 * @param theOutcomeFactory A factory object that can be used to create a new success or failure
 * 	object
 * 	
 * @param theContext The login context. This object contains details about the authorized
 * 	scopes and claims. Because this script will be used with the SMART Inbound Security
 * 	module, the type for this parameter will be of type SecurityInSmartAuthenticationContext,
 * 	which is described here: 
 * 	@see https://try.smilecdr.com/docs/javascript_execution_environment/callback_models.html#securityinsmartauthenticationcontext
 * 	
 * @returns {*} Either a successful outcome, or a failure outcome
 */
function onAuthenticateSuccess(theOutcome, theOutcomeFactory, theContext) {

    Log.info(" * Inbound Security Module: onAuthenticateSuccess() called");

	return theOutcome;
}

function onTokenGenerating(theUserSession, theAuthorizationRequestDetails) {
    Log.info(" * Inbound Security Module: onTokenGenerating() called");
}