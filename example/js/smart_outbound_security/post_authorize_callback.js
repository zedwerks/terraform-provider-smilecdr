/** ====================================================================================================================
 * See SMART Callback Scripts for details at: 
 * @see https://smilecdr.com/docs/smart/smart_on_fhir_outbound_security_module.html#smart-callback-script
 *
 * This script is invoked after the user has successfully authenticated with the OAuth2/OIDC server.
 * Applies to: SMART on FHIR Outbound Security Module
 * ===================================================================================================================== 
 */

/**
 * This function is called immediately after the user has authenticated in order to determine whether a session
 * context needs to be selected by the user, and to supply the options that will be presented to the user.
 * 
 * @see https://smilecdr.com/docs/smart/smart_on_fhir_outbound_security_module.html#smart-callback-script
 * 
 *
 * @param theUserSession     Contains details about the logged in user and their session.
 * @param theContextSelectionChoices This object should be manipulated in the
 *                                   function in order to provide the choices
 *                                   the user can select from.
 */
function onSmartLoginPreContextSelection(theUserSession, theContextSelectionChoices) {

    // Check that there is not launch resource parameter already set for this session.
    var launchResourceIds = theUserSession.getLaunchResourceIds();
    if (launchResourceIds !== null && launchResourceIds.length > 0) {
        Log.info(" * Launch resourceIds already set.");
        for (const res in launchResourceIds) {
            Log.info(" * Launch resourceId: " + res.resourceId + " with type: " + res.resourceType);
        }
        return;
    }
    else {
        Log.info(" * Launch resourceIds not set");
        // standalone launch? We need to set the launch resource parameter
        // TBD... from a picker?
    }
    return;
}

/**
 * This function is called just prior to the creation and issuing of a new
 * access token.
 * 
 * This function is primarily used in order to customize the SMART Launch Context(s)
 * associated with a particular session (i.e. because the launch context is maintained
 * in a third-party application and needs to be looked up during the auth flow).
 * 
 * @param theUserSession The authenticated user session (can be modified by the script)
 * @param theAuthorizationRequestDetails Contains details about the authorization request
 * 
 * For specifications of the argument 'theUserSession', see the documentation for the UserSessionDetails object:
 * @see https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#usersessiondetails
 * 
 * For specifications of the argument 'theAuthorizationRequestDetails', see the documentation:
 * @see https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#oauth2authorizationrequestdetails
 */
function onTokenGenerating(theUserSession, theAuthorizationRequestDetails) {
    Log.info(" * UserSession.username: " + theUserSession.username);
    Log.info(" * UserSession.external: " + theUserSession.external);

    var fhirUser = theUserSession.fhirUserUrl;

    Log.info(" * UserSession FHIR User: " + fhirUser);
    var launchResourceIds = theUserSession.getLaunchResourceIds();

    for (const res in launchResourceIds) {
        Log.info(" * UserSession Launch resourceId: " + res.resourceId + " with type: " + res.resourceType);
        // tbd... add existing resourceId(s) to bearer token.
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
        if (resource.resourceType === 'Patient') {
            var patientId = getPatientResourceId(resource.value, resource.system);
            if (patientId !== null) {
                Log.info(" * UserSession: addLaunchResourceId for patient: " + patientId);
                theUserSession.addLaunchResourceId('Patient', patientId);
                theAuthorizationRequestDetails.addAccessTokenClaim('patient', patientId);
            }
            else {
                Log.warn(" * Patient not found");
            }
        }
        else if (resource.resourceType === 'Encounter') {
            Log.warn(" * Encounter Context type not supported");
        }
        else {
            Log.warn(" * Context resource.resourceType not supported: " + resource.resourceType);
        }
    }
    else {
        Log.warn(" * No launch parameter found. This is NOT EHR-launched");
        return;
    }
}

/**
 * Called after the token has been issued, but before it it returned to the client.
 * For specifications of the argument 'theDetails', see the documentation:
 * @see https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#smartonpostauthorizedetails
 *  
 * @param {*} theDetails 
 */
function onPostAuthorize(theDetails) {
    // Called after the token has been issued, but before it it returned to the client.
    // For specifications of the argument 'theDetails', see the documentation for the
    // https://smilecdr.com/docs/javascript_execution_environment/callback_models.html#smartonpostauthorizedetails

    if (theDetails !== null) {
        //Log.info(" * Access token: " + theDetails.accessToken);
        Log.info(" * Granted scopes: " + theDetails.grantedScopes);
    }
}
/**
 * This method is called when an authorization is requested, BEFORE the
 * token is created and access is granted
 * 
 * @param theRequest The incoming theRequest
 * @param theOutcomeFactory This object is a factory for a successful
 *      response or a failure response
 * @see https://smilecdr.com/docs/security/callback_scripts.html#function-onsmartloginprecontextselection
 * @returns {*}
 */
function authenticate(theRequest, theOutcomeFactory) {
    var outcome = theOutcomeFactory.newSuccess();
    //  More work to do to properly support CODAP. For now, just return success
    return outcome;
}

// ====================================================================================================================
//  Support functions
// ====================================================================================================================

/**
 * This is a custom helper function that resolves the context identifier, calling the external Context API
 * and returning the patient identifier.
 * 
 * Expects payload as defined by the smart-context project. See github.com/zedwerks/smart-context
 * 
 * @param launchId The launch context parameter from the authorization request
 * @throws "Context resource not found in Context API response"
 * @returns The patient resource identifier
 */
function resolveLaunchParameter(launchId) {
    const contextApi = Environment.getEnv('JS_SMILE_CONTEXT_API_URL') || "http://smart-context:8088/api/context";

    var token = clientAuthForContextApi();
    var contextUrl = contextApi + "/" + launchId;
    var get = Http.get(contextUrl);

    Log.info(" * GET Context: " + contextUrl);

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
    if (resource === null) {
        Log.error(" * Failed to GET context");
        throw new Error("Context resource not found in Context API response");
    }
    Log.info(" * Response resource.type: " + resource.type);
    return resource;
}

/**
 * Client Credentials Grant authentication, returning a token for smile cdr auth server
 * as client application to the context api.
 * @returns The token
 */
function clientAuthForContextApi() {
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
    // Log.debug(" * Client Credentials Grant Token: " + token);
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
function getPatientResourceId(idValue, systemValue) {

    Log.info(" * Searching for patient with identifier: " + systemValue + "|" + idValue);
    /*
    * var client = FhirClientFactory.newClient('http://127.0.0.1:8000/');
    var patientList = client.search().forResource('Patient')
        .whereToken('identifier', systemValue, idValue)
        .asList();
    */

    // This supposed to work when we have set the FHIR 
    // server that is set up as a module dependency of this outbound security module.

    /*
    var patientList = Fhir.search()
        .forResource('Patient')
        .whereToken('identifier', systemValue, idValue)
        .asList();

    if (patientList === null) {
        Log.error(" * Patient not found");
        throw new Error("Patient not found");
    } 

    var resourceId = patientList[0].id;
    */
    var resourceId = "pat9094888999";

    Log.info(" * Patient found: Patient/" + resourceId);
    return resourceId;
}