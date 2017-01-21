using UnityEditor;
using UnityEngine;

namespace Plugins.Editor
{
    public class Build
    {
        [MenuItem("Build/Build Game and Server &x")]
        public static void BuildGameAndServer()
        {
            Debug.Log("Building the game client...");
            var originalName = PlayerSettings.productName;

            try
            {
                PlayerSettings.productName = originalName + "-client";
                var options = new BuildPlayerOptions();
                options.scenes = new[] {"Assets/Scenes/SoccerField.unity"};
                options.locationPathName = "Builds/Linux/Game.x86_64";
                options.target = BuildTarget.StandaloneLinux64;
                options.options = BuildOptions.None;
                BuildPipeline.BuildPlayer(options);

                Debug.Log("Building server...");

                PlayerSettings.productName = originalName + "-server";
                options.locationPathName = "Builds/Linux/Server.x86_64";
                options.options = BuildOptions.EnableHeadlessMode;
                BuildPipeline.BuildPlayer(options);
                Debug.Log("Building Complete!");
            }
            finally
            {
                PlayerSettings.productName = originalName;
            }
        }
    }
}