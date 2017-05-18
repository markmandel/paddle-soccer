namespace Server.Controllers
{
    /// <summary>
    /// Interface for the unity score controller
    /// </summary>
    public interface IScoreController
    {
        /// <summary>
        /// Stop the server after a number of seconds
        /// </summary>
        /// <param name="seconds"></param>
        void DelayedStopServer(int seconds);
    }
}