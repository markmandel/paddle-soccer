using Client.GUI;
using Game;
using NUnit.Framework;
using UnityEngine;
using UnityEngine.UI;

namespace Tests.Editor.Client.GUI
{
    [TestFixture]
    public class ScoreBoardLogicTests
    {
        private Text yourScore;
        private Text opponentScore;
        private Text centerMessage;

        [SetUp]
        public void Setup()
        {
            yourScore = new GameObject().AddComponent<Text>();
            opponentScore = new GameObject().AddComponent<Text>();
            centerMessage = new GameObject().AddComponent<Text>();
        }

        [Test]
        public void OnScoreNotWonOrLost()
        {
            var logic = new ScoreBoardLogic(yourScore, opponentScore, centerMessage);
            for (var i = 0; i <= 1; i++)
            {
                for (var y = 0; y < PlayerScore.WinningScore; y++)
                {
                    var isPlayerLocal = i == 0;
                    logic.OnScore(y, isPlayerLocal);
                    var check = isPlayerLocal ? yourScore : opponentScore;

                    Assert.AreEqual(centerMessage.text, "");
                    StringAssert.Contains(y.ToString(), check.text);
                }
            }
        }

        [Test]
        public void OnScoreWon()
        {
            var logic = new ScoreBoardLogic(yourScore, opponentScore, centerMessage);
            logic.OnScore(PlayerScore.WinningScore, true);

            StringAssert.Contains(PlayerScore.WinningScore.ToString(), yourScore.text);
            StringAssert.Contains("WIN", centerMessage.text);
        }
        
        [Test]
        public void OnScoreLost()
        {
            var logic = new ScoreBoardLogic(yourScore, opponentScore, centerMessage);
            logic.OnScore(PlayerScore.WinningScore, false);

            StringAssert.Contains(PlayerScore.WinningScore.ToString(), opponentScore.text);
            StringAssert.Contains("LOSE", centerMessage.text);
        }
    }
}