namespace Network
{
    public interface IUnityServer
    {
        /// <summary>
        /// Start the listening server.
        /// </summary>
        /// <returns>false if something has gone wrong</returns>
        bool StartServer();
    }
}