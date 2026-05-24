package state

import (
	"fmt"
	"path/filepath"
)

func ProviderPubspecDependency() string {
	return "  provider: ^6.1.2"
}

func GenerateProviderFiles(projectPath, featureName string) error {
	pascal := toPascalCase(featureName)
	base := filepath.Join(projectPath, fmt.Sprintf("lib/features/%s/presentation/providers", featureName))

	path := filepath.Join(base, fmt.Sprintf("%s_provider.dart", featureName))
	return writeFile(path, providerDart(pascal, featureName))
}

func providerDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/foundation.dart';

class %sProvider extends ChangeNotifier {
  bool _isLoading = false;
  List<dynamic> _data = [];
  String? _error;

  bool get isLoading => _isLoading;
  List<dynamic> get data => _data;
  String? get error => _error;
  bool get isEmpty => _data.isEmpty;

  Future<void> load() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // TODO: implement load logic
      _data = [];
    } catch (e) {
      _error = e.toString();
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }
}
`, pascal)
}
