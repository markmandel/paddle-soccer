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

        void PostHTTP(string host, string body);
    }
}