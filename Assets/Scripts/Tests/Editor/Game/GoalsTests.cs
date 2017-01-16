using Game;
using NUnit.Framework;
using UnityEngine;

namespace Tests.Editor.Game
{
    public class GoalsTests
    {
        [Test]
        public void OnBallGoal()
        {
            var check = false;
            var action = Goals.OnBallGoal(_ => check = true);

            var go = new GameObject("Foo");
            var collider = go.AddComponent<SphereCollider>();
            action(collider);
            Assert.IsFalse(check);

            collider.name = Ball.Name;
            action(collider);
            Assert.IsTrue(check);

            Object.DestroyImmediate(go);
        }
    }
}