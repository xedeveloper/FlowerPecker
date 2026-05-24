package project

import (
	"fmt"
	"path/filepath"
	"strings"
)

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	result := ""
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func GenerateMVCArchitecture(projectPath, projectName string) error {
	if err := generateCoreFiles(projectPath, projectName); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(projectPath, "lib/models"),
		filepath.Join(projectPath, "lib/views"),
		filepath.Join(projectPath, "lib/controllers"),
	}
	for _, d := range dirs {
		if err := writeGitkeep(d); err != nil {
			return err
		}
	}

	mainDart := generateMVCMainDart(projectName)
	return writeFile(filepath.Join(projectPath, "lib/main.dart"), mainDart)
}

func generateMVCMainDart(projectName string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';
import 'core/theme/app_theme.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const %sApp());
}

class %sApp extends StatelessWidget {
  const %sApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: '%s',
      theme: AppTheme.lightTheme,
      debugShowCheckedModeBanner: false,
      home: const Scaffold(
        body: Center(
          child: Text('Welcome to %s'),
        ),
      ),
    );
  }
}
`, toPascalCase(projectName), toPascalCase(projectName), toPascalCase(projectName), projectName, projectName)
}
