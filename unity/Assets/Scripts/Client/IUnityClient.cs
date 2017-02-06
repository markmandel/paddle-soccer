using System;
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