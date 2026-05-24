package project

import (
	"fmt"
	"path/filepath"
)

func GenerateCleanArchitecture(projectPath, projectName string) error {
	if err := generateCoreFiles(projectPath, projectName); err != nil {
		return err
	}

	if err := writeGitkeep(filepath.Join(projectPath, "lib/features")); err != nil {
		return err
	}

	mainDart := generateCleanMainDart(projectName)
	return writeFile(filepath.Join(projectPath, "lib/main.dart"), mainDart)
}

func generateCleanMainDart(projectName string) string {
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
