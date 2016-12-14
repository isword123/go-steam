using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.IO;
using SteamLanguageParser;

namespace GoSteamLanguageGenerator
{
	class MainClass
	{
		public static void Main(string[] args)
		{
			if (args.Length < 2) {
				Console.WriteLine("Must have at least two parameters: SteamLanguage files path and output path!");
				return;
			}

			string languagePath = Path.GetFullPath(args[0]);
			string outputPath = Path.GetFullPath(args[1]);

			Environment.CurrentDirectory = languagePath;

			var codeGen = new GoGen();

			Queue<Token> tokenList = LanguageParser.TokenizeString(File.ReadAllText("steammsg.steamd"));

			Node root = TokenAnalyzer.Analyze(tokenList);

			Node rootEnumNode = new Node();
			Node rootMessageNode = new Node();

			rootEnumNode.childNodes.AddRange(root.childNodes.Where(n => n is EnumNode));
			rootMessageNode.childNodes.AddRange(root.childNodes.Where(n => n is ClassNode));

			StringBuilder enumBuilder = new StringBuilder();
			StringBuilder messageBuilder = new StringBuilder();

			codeGen.EmitEnums(rootEnumNode, enumBuilder);
			codeGen.EmitClasses(rootMessageNode, messageBuilder);

			string outputEnumFile = Path.Combine(outputPath, "enums.go");
			string outputMessageFile = Path.Combine(outputPath, "messages.go");

            Directory.CreateDirectory(Path.GetDirectoryName(outputEnumFile));

			File.WriteAllText(outputEnumFile, enumBuilder.ToString());
			File.WriteAllText(outputMessageFile, messageBuilder.ToString());
		}
	}
}
