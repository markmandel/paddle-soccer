using UnityEngine.Networking;

namespace Client
{
    /// <summary>
    /// Interface for a network client
    /// </summary>
    public interface IUnityClient
    {
        /// <summary>
        /// Starts the client connection
        /// </summary>
        /// <returns>The Network client</returns>
        NetworkClient StartClient();

        /// <summary>
        /// Change the Server Host settings
        /// </summary>
        /// <param name="host">The server host</param>
        /// <param name="port">The port to use</param>
        void SetHost(string host, int port);
    }
}