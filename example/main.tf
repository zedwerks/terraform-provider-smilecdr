# Copyright (c) Zed Werks Inc.
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    smilecdr = {
      source  = "zedwerks/smilecdr"
      version = "0.2.0"
    }
  }
}

provider "smilecdr" {
  url      = "http://localhost:8000"
  username = "admin"
  password = "password"
}


resource "smilecdr_openid_client" "client1" {

  nodeId                     = "Master"
  moduleId                   = "smart_auth"
  accessTokenValiditySeconds = 300
  allowedGrantTypes = [
    "authorization_code",
    "refresh_token",
    "client_credentials",
    "password",
    "implicit"
  ]
  autoApproveScopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  autoGrantScopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  clientId                    = "client1"
  clientName                  = "Client 1"
  fixedScope                  = false
  refreshTokenValiditySeconds = 86400
  registeredRedirectUris      = ["http://localhost:3000"]
  scopes = [
    "openid",
    "profile",
    "email",
    "patient/*.*",
    "launch",
    "launch/patient",
    "launch/encounter",
    "offline_access"
  ]
  secretRequired         = false
  secretClientCanChange  = false
  enabled                = true
  canIntrospectAnyTokens = false
  canIntrospectOwnTokens = false
  alwaysRequireApproval  = false
  canReissueTokens       = false
  rememberApprovedScopes = false
  attestationAccepted    = false
  publicJwks             = ""
  jwksUrl                = ""
  createdByAppSphere     = false

}
