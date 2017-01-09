using UnityEditor;
using UnityEngine;

namespace Plugins.Editor
{
    public class Build
    {
        [MenuItem("Tools/Build Game and Server &x")]
        public static void BuildGameAndServer()
        {
            Debug.Log("Building the game client...");

            var options = new BuildPlayerOptions();
            options.scenes = new string[] {"Assets/Scenes/SoccerField.unity"};
            options.locationPathName = "Builds/Linux/Game.x86_64";
            options.target = BuildTarget.StandaloneLinux64;
            options.options = BuildOptions.None;

            BuildPipeline.BuildPlayer(options);

            options.locationPathName = "Builds/Linux/Server.x86_64";
            options.options = BuildOptions.EnableHeadlessMode;

            Debug.Log("Building server...");

            BuildPipeline.BuildPlayer(options);

            Debug.Log("Building Complete!");
        }
    }
}