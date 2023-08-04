// federationUserMappingScript.js
// An optional script that is used to create Smile CDR user name from the federated login details.

/**
 * This is a sample user name mapping callback script
 *
 * @param theOidcUserInfoMap OIDC claims from the token as a map
 * 
 * @param theServerInfo JSON mapping of the OAuth server definition (backed by ca.cdr.api.model.json.OAuth2ServerJson)
 * 
 * @returns Local unique Smile CDR user name for the enternal user.  
 */
function getUserName(theOidcUserInfoMap, theServerInfo) {
    return "EXT_USER:" + theOidcUserInfoMap['preferred_username'];
 }