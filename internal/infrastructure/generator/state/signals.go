package state

import (
	"fmt"
	"path/filepath"
)

func SignalsPubspecDependency() string {
	return "  signals_flutter: ^5.4.0"
}

func GenerateSignalsFiles(projectPath, featureName string) error {
	pascal := toPascalCase(featureName)
	base := filepath.Join(projectPath, fmt.Sprintf("lib/features/%s/presentation/signals", featureName))

	path := filepath.Join(base, fmt.Sprintf("%s_signals.dart", featureName))
	return writeFile(path, signalsDart(pascal, featureName))
}

func signalsDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:signals_flutter/signals_flutter.dart';

final %sIsLoading = signal<bool>(false);
final %sData = signal<List<dynamic>>([]);
final %sError = signal<String?>( null);

final %sIsEmpty = computed(() => %sData.value.isEmpty);

Future<void> load%s() async {
  %sIsLoading.value = true;
  %sError.value = null;
  try {
    // TODO: implement load logic
    %sData.value = [];
  } catch (e) {
    %sError.value = e.toString();
  } finally {
    %sIsLoading.value = false;
  }
}
`, feature, feature, feature, feature, feature, pascal, feature, feature, feature, feature, feature)
}
