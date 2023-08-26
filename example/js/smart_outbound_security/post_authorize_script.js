// post_authorize_script.js
// See SMART Callback Scripts for details at: 
// https://smilecdr.com/docs/smart/smart_on_fhir_outbound_security_module.html#smart-callback-script
// 
// This script is invoked after the user has successfully authenticated with the OAuth2/OIDC server.
// Applies to: SMART on FHIR Outbound Security Module
// 

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
function onTokenGenerating(theUserSession, theAuthorizationRequestDetails) {
    Log.info(" * UserSession.username: " + theUserSession.username);
    Log.info(" * UserSession.external: " + theUserSession.external);

    var fhirUser = theUserSession.fhirUserUrl;

    Log.info(" * UserSession FHIR User: " + fhirUser);
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
        if (resource.type === 'patient') {
            var patient = getPatientResource(resource.value, resource.system);
            if (patient !== null) {
                theUserSession.addLaunchResourceId('Patient', patient.identifier);
            }
        }
        else if (resource.type === 'encounter') {
            Log.warn(" * Encounter Context type not supported");
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
function onPostAuthorize(theDetails) {
    // Called after the token has been issued, but before it it returned to the client.
    // For specifications of the argument 'theDetails', see the documentation for the
    // https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#smartonpostauthorizedetails

    if (theDetails !== null) {
        Log.info(" * Access token: " + theDetails.accessToken);
        Log.info(" * Authorized scopes: " + theDetails.grantedScopes);
    }
}

/**
 * This is a custom helper function that resolves the context identifier, calling the external Context API
 * and returning the patient identifier.
 * 
 * Expects payload as defined by the smart-context project. See github.com/zedwerks/smart-context
 * 
 * @param launchId The launch context parameter from the authorization request
 * @returns The patient resource identifier
 */
function resolveLaunchParameter(launchId) {
    const contextApi = Environment.getEnv('JS_SMILE_CONTEXT_API_URL') || "http://smart-context:8088/api/context";

    var token = authenticate();
    var contextUrl = contextApi + "/" + launchId;
    var get = Http.get(contextUrl);

    Log.info(" * GET Context: " + contextUrl);
    Log.debug(" * Token: " + token);

    get.addRequestHeader('Accept', 'application/json');
    get.addRequestHeader('Authorization', 'Bearer ' + token);

    get.execute();
    if (!get.isSuccess()) {
        Log.error(" * Failed to GET context");
        Log.error(" * Response: " + get.getFailureMessage());
        throw get.getFailureMessage();
    }
    var responseJson = get.parseResponseAsJson();
    Log.info(" * Context Response.resourceType: " + responseJson.resourceType);
    var resource = responseJson.parameter[0].resource;
    Log.info(" * Resource: " + resource);
    return resource;
}

/**
 * Client Credentials Grant authentication, returning a token for smile cdr auth server
 * as client application to the context api.
 * @returns The token
 */
function authenticate() {
    const clientId = Environment.getEnv('JS_SMILE_CONTEXT_API_CLIENT') || "smile-cdr";
    const clientSecret = Environment.getEnv('JS_SMILE_CONTEXT_API_CLIENT_SECRET');
    const tokenEndpoint = Environment.getEnv('JS_SMILE_CONTEXT_API_TOKEN_URL');
    const scope = Environment.getEnv('JS_SMILE_CONTEXT_API_SCOPE') || "context";

    if (clientSecret === null) {
        Log.warn(" * Client Credentials Grant authentication failed");
        Log.warn(" * Client secret not set");
        return null;
    }
    if (tokenEndpoint === null) {
        Log.warn(" * Client Credentials Grant authentication failed");
        Log.warn(" * Token endpoint not set");
        return null;
    }
    Log.info(" * Client Credentials Grant Token Url: " + tokenEndpoint);
    Log.info(" * Client Credentials Grant Client Id: " + clientId);
    var post = Http.post(tokenEndpoint);
    post.setContentType('application/x-www-form-urlencoded');
    post.addRequestHeader('Accept', 'application/json');
    var clientIdEncoded = Converter.urlEncodeString(clientId);
    var clientSecretEncoded = Converter.urlEncodeString(clientSecret);
    var scopeEncoded = Converter.urlEncodeString(scope);
    var formContent = "client_id=" + clientIdEncoded + "&client_secret=" + clientSecretEncoded + "&grant_type=client_credentials&scope=" + scopeEncoded;

    post.setContentString(formContent);
    post.execute();

    if (!post.isSuccess()) {
        Log.warn(" * Client Credentials Grant authentication failed");
        Log.warn(" * Response: " + post.getFailureMessage());
        return null;
    }
    var responseJson = post.parseResponseAsJson();
    Log.info(" * Client Credentials Grant authentication");
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
function getPatientResource(idValue, systemValue) {
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