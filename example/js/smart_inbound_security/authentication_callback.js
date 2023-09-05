/** ====================================================================================================================
 * See SMART Inbound Module Authentication Callback Scripts for details at: 
 * @see https://smilecdr.com/docs/security/callback_scripts.html
 *
 * This script is invoked after the user has successfully authenticated with the OAuth2/OIDC server.
 * Applies to: SMART on FHIR Outbound Security Module
 * ===================================================================================================================== 
 */


/**
 * This is a sample authentication callback script for the
 * SMART Inbound Security module.
 *
 * @param theOutcome The outcome object. This contains details about the user that was created
 *    in response to the incoming token.
 *
 * @param theOutcomeFactory A factory object that can be used to create a new success or failure
 *    object
 *
 * @param theContext The login context. This object contains details about the authorized
 *    scopes and claims. Because this script will be used with the SMART Inbound Security
 *    module, the type for this parameter will be of type SecurityInSmartAuthenticationContext,
 *    which is described here:
 * 
 *    @see https://smilecdr.com/docs/security/callback_scripts.html
 *
 * @returns {*} Either a successful outcome, or a failure outcome
 */
function onAuthenticateSuccess(theOutcome, theOutcomeFactory, theContext) {

	Log.info(" * Inbound Security Module: onAuthenticateSuccess() called");

	setPractitionerRoleIfAny(theContext, theOutcome);
	setPatientAuthority(theContext, theOutcome);

	return theOutcome;
}

function setPatientAuthority(theContext, theOutcome) {
	var approvedScopes = theContext.getApprovedScopes();

	let patientClaim = theContext.getClaim('patient');
	let patientUri = 'Patient/' + patientClaim;

	if (patientClaim != null) {
		theOutcome.addAuthority('FHIR_CAPABILITIES');
	}

	if (patientClaim != null  && theContext.hasApprovedScope('patient/*.read')) {
		theOutcome.addAuthority('FHIR_READ_ALL_IN_COMPARTMENT', patientUri);
	}
	if (patientClaim != null && theContext.hasApprovedScope('patient/*.write')) {
		theOutcome.addAuthority('FHIR_WRITE_ALL_IN_COMPARTMENT' + patientUri);
	}
	if (patientClaim != null && theContext.hasApprovedScope('patient/*.*')) {
		theOutcome.addAuthority('FHIR_READ_ALL_IN_COMPARTMENT', patientUri);
		theOutcome.addAuthority('FHIR_READ_ALL_IN_COMPARTMENT', patientUri);
	}

	if (theContext.getClaim('patientRepresentative') != null) {}
}

function setPractitionerRoleIfAny(theContext, theOutcome) {
	// check if the user is a practitioner
	var practitionerClaim = theContext.getClaim('practitioner');
}
