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
    }
}