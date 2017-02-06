using System;
using UnityEngine.Networking;

namespace Server
{
    public interface IUnityServer
    {
        /// <summary>
        /// Start the listening server.
        /// </summary>
        /// <returns>false if something has gone wrong</returns>
        bool StartServer();

        /// <summary>
        /// Change the port the server starts on.
        /// </summary>
        /// <param name="port">The port to use</param>
        void SetPort(int port);

        /// <summary>
        /// Send a POST HTTP Request, and call the lambda when complete
        /// </summary>
        /// <param name="host">the host</param>
        /// <param name="body">the body to send (probably json)</param>
        /// <param name="lambda">optional lambda to call</param>
        void PostHTTP(string host, string body, Action<UnityWebRequest> lambda = null);
    }
}