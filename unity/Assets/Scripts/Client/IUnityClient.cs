// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

﻿using System;
using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Interface for a network client
    /// </summary>
    public interface IUnityClient
    {
        /// <summary>
        /// Change where the server host location
        /// </summary>
        /// <param name="host">The server host</param>
        void SetHost(string host);

        /// <summary>
        /// Change the port the client connect to.
        /// </summary>
        /// <param name="port">The port to use</param>
        void SetPort(int port);

        /// <summary>
        /// Starts the client connection
        /// </summary>
        /// <returns>The Network client</returns>
        NetworkClient StartClient();

        /// <summary>
        /// Send a POST HTTP Request, and call Action when complete
        /// </summary>
        /// <param name="host">the host</param>
        /// <param name="body">the body to send (probably json)</param>
        /// <param name="lambda">optional lambda to call</param>
        void PostHTTP(string host, string body, Action<UnityWebRequest> lambda = null);

        /// <summary>
        /// Polls a GET service until the lambda returns true
        /// </summary>
        /// <param name="host">The host url to call</param>
        /// <param name="lambda">The lambda to call. Will continue to poll until true is returned</param>
        void PollGetHTTP(string host, Func<UnityWebRequest, bool> lambda);
    }
}