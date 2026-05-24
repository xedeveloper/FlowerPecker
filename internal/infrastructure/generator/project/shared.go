package project

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func writeGitkeep(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, ".gitkeep"), []byte(""), 0644)
}

func appThemeDart() string {
	return `import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppColors {
  static const Color primary = Color(0xFF000000);
  static const Color onPrimary = Color(0xFFFFFFFF);
  static const Color ink = Color(0xFF000000);
  static const Color inkSoft = Color(0xFF1A1A1A);
  static const Color body = Color(0xFF757575);
  static const Color hairline = Color(0xFFE0E0E0);
  static const Color canvas = Color(0xFFFFFFFF);
  static const Color canvasSoft = Color(0xFFF5F5F5);
  static const Color link = Color(0xFF057DBC);
}

class AppTextStyles {
  static TextStyle get displayHero => GoogleFonts.playfairDisplay(
        fontSize: 64,
        fontWeight: FontWeight.w400,
        height: 59.52 / 64,
        letterSpacing: -0.5,
        color: AppColors.ink,
      );

  static TextStyle get displayLg => GoogleFonts.playfairDisplay(
        fontSize: 48,
        fontWeight: FontWeight.w400,
        height: 1.1,
        letterSpacing: -0.3,
        color: AppColors.ink,
      );

  static TextStyle get displayMd => GoogleFonts.playfairDisplay(
        fontSize: 36,
        fontWeight: FontWeight.w400,
        height: 1.15,
        color: AppColors.ink,
      );

  static TextStyle get displaySm => GoogleFonts.playfairDisplay(
        fontSize: 28,
        fontWeight: FontWeight.w400,
        height: 1.2,
        color: AppColors.ink,
      );

  static TextStyle get displayXs => GoogleFonts.playfairDisplay(
        fontSize: 22,
        fontWeight: FontWeight.w400,
        height: 1.25,
        color: AppColors.ink,
      );

  static TextStyle get bodyLg => GoogleFonts.lora(
        fontSize: 18,
        fontWeight: FontWeight.w400,
        height: 1.6,
        color: AppColors.inkSoft,
      );

  static TextStyle get bodyMd => GoogleFonts.lora(
        fontSize: 16,
        fontWeight: FontWeight.w400,
        height: 1.6,
        color: AppColors.inkSoft,
      );

  static TextStyle get bodySm => GoogleFonts.lora(
        fontSize: 14,
        fontWeight: FontWeight.w400,
        height: 1.5,
        color: AppColors.body,
      );

  static TextStyle get bodySmStrong => GoogleFonts.lora(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        height: 1.5,
        color: AppColors.inkSoft,
      );

  static TextStyle get labelMd => const TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.5,
        color: AppColors.ink,
      );

  static TextStyle get labelSm => const TextStyle(
        fontSize: 12,
        fontWeight: FontWeight.w600,
        letterSpacing: 0.8,
        color: AppColors.body,
      );
}

class AppTheme {
  static ThemeData get lightTheme => ThemeData(
        useMaterial3: true,
        colorScheme: const ColorScheme.light(
          primary: AppColors.primary,
          onPrimary: AppColors.onPrimary,
          surface: AppColors.canvas,
          onSurface: AppColors.ink,
        ),
        scaffoldBackgroundColor: AppColors.canvas,
        appBarTheme: const AppBarTheme(
          backgroundColor: AppColors.canvas,
          foregroundColor: AppColors.ink,
          elevation: 0,
          centerTitle: true,
        ),
        textTheme: TextTheme(
          displayLarge: AppTextStyles.displayHero,
          displayMedium: AppTextStyles.displayLg,
          displaySmall: AppTextStyles.displayMd,
          headlineMedium: AppTextStyles.displaySm,
          headlineSmall: AppTextStyles.displayXs,
          bodyLarge: AppTextStyles.bodyLg,
          bodyMedium: AppTextStyles.bodyMd,
          bodySmall: AppTextStyles.bodySm,
          labelMedium: AppTextStyles.labelMd,
          labelSmall: AppTextStyles.labelSm,
        ),
        dividerTheme: const DividerThemeData(
          color: AppColors.hairline,
          thickness: 1,
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.primary,
            foregroundColor: AppColors.onPrimary,
            shape: const RoundedRectangleBorder(),
            elevation: 0,
          ),
        ),
        outlinedButtonTheme: OutlinedButtonThemeData(
          style: OutlinedButton.styleFrom(
            foregroundColor: AppColors.ink,
            side: const BorderSide(color: AppColors.ink),
            shape: const RoundedRectangleBorder(),
          ),
        ),
        inputDecorationTheme: const InputDecorationTheme(
          border: OutlineInputBorder(
            borderRadius: BorderRadius.zero,
            borderSide: BorderSide(color: AppColors.ink),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.zero,
            borderSide: BorderSide(color: AppColors.ink, width: 2),
          ),
        ),
      );
}
`
}

func appConstantsDart(projectName string) string {
	return fmt.Sprintf(`class AppConstants {
  static const String appName = '%s';
  static const String appVersion = '1.0.0';
}
`, projectName)
}

func failuresDart() string {
	return `abstract class Failure {
  final String message;
  const Failure(this.message);
}

class ServerFailure extends Failure {
  const ServerFailure(super.message);
}

class CacheFailure extends Failure {
  const CacheFailure(super.message);
}

class NetworkFailure extends Failure {
  const NetworkFailure(super.message);
}

class ValidationFailure extends Failure {
  const ValidationFailure(super.message);
}
`
}

func extensionsDart() string {
	return `import 'package:flutter/material.dart';

extension ContextExtensions on BuildContext {
  ThemeData get theme => Theme.of(this);
  TextTheme get textTheme => Theme.of(this).textTheme;
  ColorScheme get colorScheme => Theme.of(this).colorScheme;
  double get screenWidth => MediaQuery.of(this).size.width;
  double get screenHeight => MediaQuery.of(this).size.height;
  bool get isMobile => screenWidth < 768;
  bool get isTablet => screenWidth >= 768 && screenWidth < 1024;
  bool get isDesktop => screenWidth >= 1024;
}

extension StringExtensions on String {
  String get capitalize => isEmpty ? this : '${this[0].toUpperCase()}${substring(1)}';
  String get titleCase => split(' ').map((w) => w.capitalize).join(' ');
  bool get isValidEmail => RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(this);
}
`
}

func generateCoreFiles(projectPath, projectName string) error {
	coreFiles := map[string]string{
		"lib/core/constants/app_constants.dart": appConstantsDart(projectName),
		"lib/core/errors/failures.dart":         failuresDart(),
		"lib/core/theme/app_theme.dart":         appThemeDart(),
		"lib/core/utils/extensions.dart":        extensionsDart(),
	}

	for relPath, content := range coreFiles {
		fullPath := filepath.Join(projectPath, relPath)
		if err := writeFile(fullPath, content); err != nil {
			return fmt.Errorf("writing %s: %w", relPath, err)
		}
	}
	return nil
}
