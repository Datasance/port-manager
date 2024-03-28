/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package client

import (
	"encoding/json"
	"bytes"
	"fmt"
)
// create user can be removed!!
func (clt *Client) CreateUser(request User) error {
	// Send request
	if _, err := clt.doRequest("POST", "/user/signup", request); err != nil {
		return err
	}

	return nil
}

func (clt *Client) Login(request LoginRequest) (err error) {

	// Prompt for OTP
	promptOTP := func() string {
		fmt.Println("Enter OTP:")
		var otp string
		fmt.Scanln(&otp)

		return otp
	}
	otp := promptOTP()
	request.Totp = otp

	// Send request
	body, err := clt.doRequest("POST", "/user/login", request)
	if err != nil {
		return
	}

	// Read access token from response
	var response LoginResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	clt.accessToken = response.AccessToken

	return
}

func (clt *Client) RefreshUserSubscriptionKeyCtl(request LoginRequest) (err error, userSubscriptionKey string ) {
	// Send request
	bodyLogin, errLogin := clt.doRequest("POST", "/user/login", request)
	if errLogin != nil {
		return errLogin, ""
	}
	
	// Read access token from response
	var response LoginResponse
	errLoginMarshal := json.Unmarshal(bodyLogin, &response);
	if errLoginMarshal != nil {
		return errLoginMarshal, "" 
	}

	clt.SetAccessToken(response.AccessToken)

	headers := map[string]string{
		"Authorization": clt.GetAccessToken(),
		"Content-Type": "application/json",
	}

	emptyBody := bytes.NewBuffer([]byte{})


	bodyGetUser, errGetUser := clt.doRequestWithHeaders("GET", "/user/profile",emptyBody , headers)

	if errGetUser != nil {
		return errGetUser, ""
	}

	var userResponse UserResponse
	errGetUserMarshal := json.Unmarshal(bodyGetUser, &userResponse);
	if errGetUserMarshal != nil {
		return errGetUserMarshal, ""
	}

	return nil, userResponse.SubscriptionKey

}

func (clt *Client) UpdateUserPassword(request UpdateUserPasswordRequest) (err error) {
	// Send request
	_, err = clt.doRequest("PATCH", "/user/password", request)
	if err != nil {
		return
	}

	return
}
