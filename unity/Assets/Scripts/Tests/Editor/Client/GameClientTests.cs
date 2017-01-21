using System;
using Client;
using NSubstitute;
using NUnit.Framework;

namespace Tests.Editor.Client
{
    [TestFixture]
    public class GameClientTests
    {
        private IUnityClient unityServer;

        [SetUp]
        public void Setup()
        {
            unityServer = Substitute.For<IUnityClient>();
        }

        [TearDown]
        public void Teardown()
        {
            GameClient.Stop();
        }

        [Test]
        public void Start()
        {
            GameClient.Start(unityServer, new string[0]);
            unityServer.Received(1).StartClient();
            Assert.Throws<Exception>(() => GameClient.Start(unityServer, new string[0]));
        }

        [Test]
        public void StartWithHost()
        {
            var host = "10.10.10.10";
            GameClient.Start(unityServer, new[] {"-host", host});
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost(host, 7777);
        }

        [Test]
        public void StartWithPort()
        {
            var args = new[] {"-port", "8080"};
            GameClient.Start(unityServer, args);
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost("localhost", 8080);
        }

        [Test]
        public void StartWithHostAndPort()
        {
            var host = "10.10.10.10";
            var args = new[] {"-host", host, "-port", "8080"};
            GameClient.Start(unityServer, args);
            unityServer.Received(1).StartClient();
            unityServer.Received(1).SetHost(host, 8080);
        }
    }
}