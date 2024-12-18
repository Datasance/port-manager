/*
 *  *******************************************************************************
 *  * Copyright (c) 2024 Datasance Teknoloji A.S.
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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	// Prompt for OTP if not already set
	if request.Totp == "" {
		fmt.Println("Enter OTP: \n")
		reader := bufio.NewReader(os.Stdin)
		otp, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read OTP: %v", err)
		}
		request.Totp = strings.TrimSpace(otp)
	}

	// Send login request
	body, err := clt.doRequest("POST", "/user/login", request)
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	// Parse response
	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}

	clt.accessToken = response.AccessToken
	clt.refreshToken = response.RefreshToken

	return nil
}

func (clt *Client) Refresh(request RefreshTokenRequest) (err error) {
	// Send refresh request
	body, err := clt.doRequest("POST", "/user/refresh", request)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}
	var response LoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse refresh response: %v", err)
	}

	clt.accessToken = response.AccessToken
	clt.refreshToken = response.RefreshToken

	return nil
}

func (clt *Client) Profile(request WithTokenRequest) (err error, userResponse UserResponse) {
	clt.SetAccessToken(request.AccessToken)

	headers := map[string]string{
		"Authorization": clt.GetAccessToken(),
		"Content-Type":  "application/json",
	}

	bodyGetUser, err := clt.doRequestWithHeaders("GET", "/user/profile", nil, headers)
	if err != nil {
		return fmt.Errorf("failed to fetch user profile: %v", err), UserResponse{}
	}

	if err := json.Unmarshal(bodyGetUser, &userResponse); err != nil {
		return fmt.Errorf("failed to parse user profile: %v", err), UserResponse{}
	}

	return nil, userResponse
}

func (clt *Client) UpdateUserPassword(request UpdateUserPasswordRequest) (err error) {
	// Send request
	_, err = clt.doRequest("PATCH", "/user/password", request)
	if err != nil {
		return
	}

	return
}
