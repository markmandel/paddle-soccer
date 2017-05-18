using Game;
using NSubstitute;
using NUnit.Framework;
using Server.Controllers;

namespace Tests.Editor.Server.Controllers
{
    [TestFixture]
    public class ScoreControllerLogicTests
    {
        private IScoreController controller;

        [SetUp]
        public void Setup()
        {
            controller = Substitute.For<IScoreController>();
        }

        [Test]
        public void DisconnectOnWin()
        {
            var logic = new ScoreControllerLogic(controller);
            for (var i = 0; i < PlayerScore.WinningScore; i++)
            {
                logic.DisconnectOnWin(i, false);
                controller.Received(0).DelayedStopServer(5);
            }
            
            logic.DisconnectOnWin(PlayerScore.WinningScore, false);
            controller.Received(1).DelayedStopServer(5);
        }
    }
}