package state

import (
	"fmt"
	"path/filepath"
)

func GetXPubspecDependency() string {
	return "  get: ^4.6.6"
}

func GenerateGetXFiles(projectPath, featureName string) error {
	pascal := toPascalCase(featureName)
	base := filepath.Join(projectPath, fmt.Sprintf("lib/features/%s/presentation", featureName))

	files := map[string]string{
		filepath.Join(base, fmt.Sprintf("controllers/%s_controller.dart", featureName)): getxControllerDart(pascal, featureName),
		filepath.Join(base, fmt.Sprintf("bindings/%s_binding.dart", featureName)):       getxBindingDart(pascal, featureName),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func getxControllerDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:get/get.dart';

class %sController extends GetxController {
  final isLoading = false.obs;
  final data = <dynamic>[].obs;
  final error = Rxn<String>();

  bool get isEmpty => data.isEmpty;

  @override
  void onInit() {
    super.onInit();
    load();
  }

  Future<void> load() async {
    isLoading.value = true;
    error.value = null;
    try {
      // TODO: implement load logic
      data.value = [];
    } catch (e) {
      error.value = e.toString();
    } finally {
      isLoading.value = false;
    }
  }
}
`, pascal)
}

func getxBindingDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:get/get.dart';
import '../controllers/%s_controller.dart';

class %sBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<  %sController>(() => %sController());
  }
}
`, feature, pascal, pascal, pascal)
}
