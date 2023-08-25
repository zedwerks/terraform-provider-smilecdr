// post_authorize_script.js
// See SMART Callback Scripts for details at: 
// https://smilecdr.com/docs/smart/smart_on_fhir_outbound_security_module.html#smart-callback-script
// 
// This script is invoked after the user has successfully authenticated with the OAuth2/OIDC server.
// Applies to: SMART on FHIR Outbound Security Module
// 

/**
 * This function is called immediately after the user has successfully authenticated in order to
 * determine whether a session context needs to be selected by the user, and if so,
 * to provide the list of available options to the user.
 * @param theUserSession The authenticated user session (can be modified by the script)
 * @param theContextSelectionChoices The list of available context selection choices (can be modified by the script)
 */
onSmartLoginPreContextSelection(theUserSession, theContextSelectionChoices) 
{
    Log.info(" * ContextSelectionChoices.haveChoices():" + theContextSelectionChoices.haveChoices());
}

/**
 * This function is called just prior to the creation and issuing of a new
 * access token.
 * 
 * @param theUserSession The authenticated user session (can be modified by the script)
 * @param theAuthorizationRequestDetails Contains details about the authorization request
 * 
 * For specifications of the argument 'theUserSession', see the documentation for the UserSessionDetails object:
 * https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#usersessiondetails
 * 
 * For specifications of the argument 'theAuthorizationRequestDetails', see the documentation:
 * https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#oauth2authorizationrequestdetails
 */
function onTokenGenerating(theUserSession, theAuthorizationRequestDetails) 
{ 
    Log.info(" * UserSession.username: " + theUserSession.username);
    Log.info(" * UserSession.external: " + theUserSession.external);
    
    var fhirContext = theUserSession.getFhirContext();

    Log.info(" * UserSession FHIR Context: " + fhirContext.reference + " with role: " + fhirContext.role);
    var launchResourceIds = theUserSession.getLaunchResourceIds();

    for (const res in launchResourceIds) {
        Log.info(" * UserSession Launch resourceId: " + res.resourceId + " with type: " + res.resourceType);
    }

    Log.info(" * Client ID: " + theAuthorizationRequestDetails.clientId);
    Log.info(" * Member ID: " + theAuthorizationRequestDetails.memberId);

    var launchParam = theAuthorizationRequestDetails.getLaunch();

    Log.info(" * Launch parameter: " + launchParam);

    // Now, if the launch parameter is present, we now need to resolve that opaque launch parameter
    // to a patient resource. This is done by calling the resolveLaunchParameter() function.
    if (launchParam !== null) {
        // sort out the patient resource id from the launch parameter
        var resource = resolveLaunchParameter(launchParam);
        if (resource.type === 'patient')
        {
            var patient = getPatientResource(resource.value, resource.system);
            theUserSession.addLaunchResourceId('Patient', patient.identifier);
        }
    }
}  

/**
 *  Called after the token has been issued, but before it it returned to the client.
 * For specifications of the argument 'theDetails', see the documentation:
 * https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#smartonpostauthorizedetails
 *  
 * @param {*} theDetails 
 */
function onPostAuthorize(theDetails) 
{
    // Called after the token has been issued, but before it it returned to the client.
    // For specifications of the argument 'theDetails', see the documentation for the
    // https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#smartonpostauthorizedetails
    
    Log.info(" * Access token: " + theDetails.accessToken);
    Log.info(" * Authorized scopes: " + theDetails.grantedScopes);
    Log.info(" * Requesting practitioner: " + theDetails.requestingPractitioner.identifier.value);
}

/**
 * This is a custom helper function that resolves the context identifier, calling the external Context API
 * and returning the patient identifier.
 * 
 * @param launchId The launch context parameter from the authorization request
 * @returns The patient resource identifier
 */
const clientId = Environment.getProperty('js.contextApi.clientId') || "smile-cdr";
const clientSecret = Environment.getProperty('js.contextApi.clientSecret') || "ck1mvyXGf1GJTSE8YNlrePIt1xDisM1N";
const tokenEndpoint = Environment.getProperty('js.contextApi.tokenEndpoint') || "http://localhost:8080/auth/realms/poc/protocol/openid-connect/token";
const scope = Environment.getProperty('js.contextApi.scope') || "launch context openid";
const contextApi = Environment.getProperty('js.contextApi.url') || "http://smart-context:8088/api/context/";

function resolveLaunchParameter(launchId)
{
    var token = authenticate();
    var get = Http.get(contextApi + launchId);
    get.addRequestHeader('Authorization', 'Bearer ' + token);
    get.execute();
    if (!get.isSuccess()) {
        throw get.getFailureMessage();
    }
    var responseJson = get.getResponseAsJson();
    Log.info(" * Response JSON: " + responseJson);
    var resource = parameter[0].resource;
    Log.info(" * Resource: " + resource);
    return resource;
}




/**
 * Client Credentials Grant authentication, returning a token for smile cdr auth server
 * as client application to the context api.
 * @returns The token
 */
function authenticate()
{
    const basicCredentials = Base64.encode(clientId + ":" + clientSecret);

    var post = Http.post(tokenEndpoint);
    post.addRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    post.addRequestHeader('Authorization', 'Basic ' + basicCredentials);
    post.setRequestBody("grant_type=client_credentials\n&scope=" + scope + "\n");

    post.execute();
    if (!post.isSuccess()) {
        throw post.getFailureMessage();
    }
    var responseJson = post.getResponseAsJson();
    Log.info(" * Client Credentials Grant authentication");
    Log.info(" * Response JSON: " + responseJson);
    var token = responseJson.access_token;
    Log.info(" * Token: " + token);
    return token;
}

/**
 * This function retrieves the patient resource from the FHIR server.
 * This can interact with the fhir server that is set up as a
 * module dependency of this outbound security module.
 * 
 * @param {*} idValue
 * @param {*} systemValue
 * @returns The patient resource
 * @throws "Patient not found"
 */
function getPatientResource(idValue, systemValue)
{
    var patientList = Fhir.search()
        .forResource('Patient')
        .whereToken('identifier', systemValue, idValue)
        .asList();

    if (patientList.length === 0) {
        Log.info(" * Patient not found");
        throw "Patient not found";
    }
    Log.info(" * Patient found: " + patientList[0].id);
    return patientList[0];
}