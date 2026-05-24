package state

import (
	"fmt"
	"path/filepath"
)

func RiverpodPubspecDependency() string {
	return "  flutter_riverpod: ^2.5.1"
}

func GenerateRiverpodFiles(projectPath, featureName string) error {
	pascal := toPascalCase(featureName)
	base := filepath.Join(projectPath, fmt.Sprintf("lib/features/%s/presentation/providers", featureName))

	files := map[string]string{
		filepath.Join(base, fmt.Sprintf("%s_provider.dart", featureName)):  riverpodProviderDart(pascal, featureName),
		filepath.Join(base, fmt.Sprintf("%s_notifier.dart", featureName)): riverpodNotifierDart(pascal, featureName),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func riverpodProviderDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter_riverpod/flutter_riverpod.dart';
import '%s_notifier.dart';

final %sProvider = StateNotifierProvider<
    %sNotifier,
    %sState>((ref) => %sNotifier());
`, feature, feature, pascal, pascal, pascal)
}

func riverpodNotifierDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter_riverpod/flutter_riverpod.dart';

class %sState {
  final bool isLoading;
  final List<dynamic> data;
  final String? error;

  const %sState({
    this.isLoading = false,
    this.data = const [],
    this.error,
  });

  %sState copyWith({
    bool? isLoading,
    List<dynamic>? data,
    String? error,
  }) {
    return %sState(
      isLoading: isLoading ?? this.isLoading,
      data: data ?? this.data,
      error: error ?? this.error,
    );
  }
}

class %sNotifier extends StateNotifier<%sState> {
  %sNotifier() : super(const %sState());

  Future<void> load() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      // TODO: implement load logic
      state = state.copyWith(isLoading: false, data: []);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}
`, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}
