package project

import (
	"fmt"
	"path/filepath"
)

func GenerateMVVMArchitecture(projectPath, projectName string) error {
	if err := generateCoreFiles(projectPath, projectName); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(projectPath, "lib/models"),
		filepath.Join(projectPath, "lib/views"),
		filepath.Join(projectPath, "lib/viewmodels"),
	}
	for _, d := range dirs {
		if err := writeGitkeep(d); err != nil {
			return err
		}
	}

	mainDart := generateMVVMMainDart(projectName)
	return writeFile(filepath.Join(projectPath, "lib/main.dart"), mainDart)
}

func generateMVVMMainDart(projectName string) string {
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
