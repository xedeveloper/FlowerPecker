package usecases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xedeveloper/flowerpecker/internal/domain/models"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/readme"
	"github.com/xedeveloper/flowerpecker/internal/infrastructure/generator/state"
)

func CreateModule(config models.ModuleConfig, onStep func(string)) error {
	onStep(fmt.Sprintf("Scaffolding %s module...", config.Name))

	switch config.Architecture {
	case models.ArchClean:
		if err := generateCleanModule(config); err != nil {
			return err
		}
	case models.ArchMVC:
		if err := generateMVCModule(config); err != nil {
			return err
		}
	case models.ArchMVVM:
		if err := generateMVVMModule(config); err != nil {
			return err
		}
	}

	onStep("Generating state management files...")
	if err := generateModuleStateFiles(config); err != nil {
		return err
	}

	onStep("Updating README...")
	projectConfig := models.ProjectConfig{
		Architecture:    config.Architecture,
		StateManagement: config.StateManagement,
	}
	return readme.UpdateReadmeWithModule(config.ProjectPath, config.Name, projectConfig)
}

func generateCleanModule(config models.ModuleConfig) error {
	base := filepath.Join(config.ProjectPath, "lib/features", config.Name)
	pascal := toPascalCase(config.Name)

		files := map[string]string{
			filepath.Join(base, "data/datasources", config.Name+"_remote_datasource.dart"): remoteDataSourceDart(pascal, config.Name),
			filepath.Join(base, "data/datasources", config.Name+"_local_datasource.dart"):  localDataSourceDart(pascal, config.Name),
			filepath.Join(base, "data/models", config.Name+"_model.dart"):                  modelDart(pascal, config.Name),
			filepath.Join(base, "data/repositories", config.Name+"_repository_impl.dart"):  repoImplDart(pascal, config.Name),
			filepath.Join(base, "domain/entities", config.Name+"_entity.dart"):             entityDart(pascal),
			filepath.Join(base, "domain/repositories", config.Name+"_repository.dart"):     abstractRepoDart(pascal, config.Name),
			filepath.Join(base, "domain/usecases", config.Name+"_usecase.dart"):            usecaseDart(pascal, config.Name),
			filepath.Join(base, "presentation/pages", config.Name+"_page.dart"):            pageDart(pascal, config.Name),
			filepath.Join(base, "presentation/widgets", config.Name+"_form_widget.dart"):   formWidgetDart(pascal, config.Name),
		}

	for path, content := range files {
		if err := writeModuleFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateMVCModule(config models.ModuleConfig) error {
	pascal := toPascalCase(config.Name)
	base := config.ProjectPath + "/lib"

	files := map[string]string{
		filepath.Join(base, "models", config.Name+"_model.dart"):          modelDart(pascal, config.Name),
		filepath.Join(base, "views", config.Name+"_view.dart"):            mvcViewDart(pascal, config.Name),
		filepath.Join(base, "controllers", config.Name+"_controller.dart"): mvcControllerDart(pascal, config.Name),
	}

	for path, content := range files {
		if err := writeModuleFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateMVVMModule(config models.ModuleConfig) error {
	pascal := toPascalCase(config.Name)
	base := config.ProjectPath + "/lib"

	files := map[string]string{
		filepath.Join(base, "models", config.Name+"_model.dart"):        modelDart(pascal, config.Name),
		filepath.Join(base, "views", config.Name+"_view.dart"):          mvvmViewDart(pascal, config.Name),
		filepath.Join(base, "viewmodels", config.Name+"_viewmodel.dart"): mvvmViewModelDart(pascal, config.Name),
	}

	for path, content := range files {
		if err := writeModuleFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

func generateModuleStateFiles(config models.ModuleConfig) error {
	switch config.StateManagement {
	case models.StateBLoC:
		return state.GenerateBlocFiles(config.ProjectPath, config.Name)
	case models.StateRiverpod:
		return state.GenerateRiverpodFiles(config.ProjectPath, config.Name)
	case models.StateSignals:
		return state.GenerateSignalsFiles(config.ProjectPath, config.Name)
	case models.StateProvider:
		return state.GenerateProviderFiles(config.ProjectPath, config.Name)
	case models.StateGetX:
		return state.GenerateGetXFiles(config.ProjectPath, config.Name)
	}
	return nil
}

func writeModuleFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func toPascalCase(s string) string {
	result := ""
	capitalizeNext := true
	for _, c := range s {
		if c == '_' {
			capitalizeNext = true
		} else if capitalizeNext {
			result += string(c - 32)
			capitalizeNext = false
		} else {
			result += string(c)
		}
	}
	return result
}

func remoteDataSourceDart(pascal, feature string) string {
	return fmt.Sprintf(`abstract class %sRemoteDataSource {
  Future<List<dynamic>> fetchAll();
}

class %sRemoteDataSourceImpl implements %sRemoteDataSource {
  @override
  Future<List<dynamic>> fetchAll() async {
    // TODO: implement remote data fetching
    return [];
  }
}
`, pascal, pascal, pascal)
}

func localDataSourceDart(pascal, feature string) string {
	return fmt.Sprintf(`abstract class %sLocalDataSource {
  Future<List<dynamic>> getCached();
  Future<void> cache(List<dynamic> data);
}

class %sLocalDataSourceImpl implements %sLocalDataSource {
  @override
  Future<List<dynamic>> getCached() async {
    // TODO: implement local cache retrieval
    return [];
  }

  @override
  Future<void> cache(List<dynamic> data) async {
    // TODO: implement local caching
  }
}
`, pascal, pascal, pascal)
}

func modelDart(pascal, feature string) string {
	return fmt.Sprintf(`class %sModel {
  final String id;
  final String name;

  const %sModel({required this.id, required this.name});

  factory %sModel.fromJson(Map<String, dynamic> json) {
    return %sModel(
      id: json['id'] as String,
      name: json['name'] as String,
    );
  }

  Map<String, dynamic> toJson() => {'id': id, 'name': name};
}
`, pascal, pascal, pascal, pascal)
}

func repoImplDart(pascal, feature string) string {
	return fmt.Sprintf(`import '../datasources/%s_remote_datasource.dart';
import '../datasources/%s_local_datasource.dart';
import '../../domain/repositories/%s_repository.dart';

class %sRepositoryImpl implements %sRepository {
  final %sRemoteDataSource remoteDataSource;
  final %sLocalDataSource localDataSource;

  %sRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<List<dynamic>> getAll() async {
    try {
      final data = await remoteDataSource.fetchAll();
      await localDataSource.cache(data);
      return data;
    } catch (_) {
      return localDataSource.getCached();
    }
  }
}
`, feature, feature, feature, pascal, pascal, pascal, pascal, pascal)
}

func entityDart(pascal string) string {
	return fmt.Sprintf(`class %sEntity {
  final String id;
  final String name;

  const %sEntity({required this.id, required this.name});
}
`, pascal, pascal)
}

func abstractRepoDart(pascal, feature string) string {
	return fmt.Sprintf(`abstract class %sRepository {
  Future<List<dynamic>> getAll();
}
`, pascal)
}

func usecaseDart(pascal, feature string) string {
	return fmt.Sprintf(`import '../repositories/%s_repository.dart';

class Get%sList {
  final %sRepository repository;

  Get%sList({required this.repository});

  Future<List<dynamic>> call() => repository.getAll();
}
`, feature, pascal, pascal, pascal)
}

func pageDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';

class %sPage extends StatelessWidget {
  const %sPage({super.key});

  static const routeName = '/%s';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('%s')),
      body: const Center(
        child: Text('%s Page'),
      ),
    );
  }
}
`, pascal, pascal, feature, pascal, pascal)
}

func formWidgetDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';

class %sFormWidget extends StatefulWidget {
  final VoidCallback? onSubmit;

  const %sFormWidget({super.key, this.onSubmit});

  @override
  State<%sFormWidget> createState() => _%sFormWidgetState();
}

class _%sFormWidgetState extends State<%sFormWidget> {
  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context) {
    return Form(
      key: _formKey,
      child: Column(
        children: [
          TextFormField(
            decoration: const InputDecoration(labelText: 'Name'),
            validator: (v) => v == null || v.isEmpty ? 'Required' : null,
          ),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () {
              if (_formKey.currentState!.validate()) {
                widget.onSubmit?.call();
              }
            },
            child: const Text('Submit'),
          ),
        ],
      ),
    );
  }
}
`, pascal, pascal, pascal, pascal, pascal, pascal)
}

func mvcViewDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';
import '../controllers/%s_controller.dart';

class %sView extends StatefulWidget {
  const %sView({super.key});

  @override
  State<%sView> createState() => _%sViewState();
}

class _%sViewState extends State<%sView> {
  late %sController _controller;

  @override
  void initState() {
    super.initState();
    _controller = %sController();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('%s')),
      body: const Center(child: Text('%s View')),
    );
  }
}
`, feature, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

func mvcControllerDart(pascal, feature string) string {
	return fmt.Sprintf(`class %sController {
  List<dynamic> items = [];
  bool isLoading = false;
  String? error;

  Future<void> loadItems() async {
    isLoading = true;
    try {
      // TODO: implement load logic
      items = [];
    } catch (e) {
      error = e.toString();
    } finally {
      isLoading = false;
    }
  }
}
`, pascal)
}

func mvvmViewDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/material.dart';
import '../viewmodels/%s_viewmodel.dart';

class %sView extends StatefulWidget {
  const %sView({super.key});

  @override
  State<%sView> createState() => _%sViewState();
}

class _%sViewState extends State<%sView> {
  late %sViewModel _viewModel;

  @override
  void initState() {
    super.initState();
    _viewModel = %sViewModel();
    _viewModel.addListener(() => setState(() {}));
    _viewModel.load();
  }

  @override
  void dispose() {
    _viewModel.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('%s')),
      body: _viewModel.isLoading
          ? const Center(child: CircularProgressIndicator())
          : const Center(child: Text('%s View')),
    );
  }
}
`, feature, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal, pascal)
}

func mvvmViewModelDart(pascal, feature string) string {
	return fmt.Sprintf(`import 'package:flutter/foundation.dart';

class %sViewModel extends ChangeNotifier {
  bool _isLoading = false;
  List<dynamic> _items = [];
  String? _error;

  bool get isLoading => _isLoading;
  List<dynamic> get items => _items;
  String? get error => _error;

  Future<void> load() async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // TODO: implement load logic
      _items = [];
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
