// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/weaveworks/reignite/core/application"
	"github.com/weaveworks/reignite/core/ports"
	"github.com/weaveworks/reignite/infrastructure/containerd"
	"github.com/weaveworks/reignite/infrastructure/controllers"
	"github.com/weaveworks/reignite/infrastructure/firecracker"
	"github.com/weaveworks/reignite/infrastructure/grpc"
	"github.com/weaveworks/reignite/infrastructure/network"
	"github.com/weaveworks/reignite/infrastructure/ulid"
	"github.com/weaveworks/reignite/internal/config"
)

// Injectors from wire.go:

func InitializePorts(cfg *config.Config) (*ports.Collection, error) {
	config2 := containerdConfig(cfg)
	microVMRepository, err := containerd.NewMicroVMRepo(config2)
	if err != nil {
		return nil, err
	}
	config3 := firecrackerConfig(cfg)
	config4 := networkConfig(cfg)
	networkService := network.New(config4)
	fs := afero.NewOsFs()
	microVMService := firecracker.New(config3, networkService, fs)
	eventService, err := containerd.NewEventService(config2)
	if err != nil {
		return nil, err
	}
	idService := ulid.New()
	imageService, err := containerd.NewImageService(config2)
	if err != nil {
		return nil, err
	}
	collection := appPorts(microVMRepository, microVMService, eventService, idService, networkService, imageService, fs)
	return collection, nil
}

func InitializeApp(cfg *config.Config, ports2 *ports.Collection) application.App {
	applicationConfig := appConfig(cfg)
	app := application.New(applicationConfig, ports2)
	return app
}

func InializeController(app application.App, ports2 *ports.Collection) *controllers.MicroVMController {
	eventService := eventSvcFromScope(ports2)
	reconcileMicroVMsUseCase := reconcileUCFromApp(app)
	microVMController := controllers.New(eventService, reconcileMicroVMsUseCase)
	return microVMController
}

func InitializeGRPCServer(app application.App) ports.MicroVMGRPCService {
	microVMCommandUseCases := commandUCFromApp(app)
	microVMQueryUseCases := queryUCFromApp(app)
	microVMGRPCService := grpc.NewServer(microVMCommandUseCases, microVMQueryUseCases)
	return microVMGRPCService
}

// wire.go:

func containerdConfig(cfg *config.Config) *containerd.Config {
	return &containerd.Config{
		SnapshotterKernel: cfg.CtrSnapshotterKernel,
		SnapshotterVolume: cfg.CtrSnapshotterVolume,
		SocketPath:        cfg.CtrSocketPath,
		Namespace:         cfg.CtrNamespace,
	}
}

func firecrackerConfig(cfg *config.Config) *firecracker.Config {
	return &firecracker.Config{
		FirecrackerBin: cfg.FirecrackerBin,
		RunDetached:    cfg.FirecrackerDetatch,
		APIConfig:      cfg.FirecrackerUseAPI,
		StateRoot:      fmt.Sprintf("%s/vm", cfg.StateRootDir),
	}
}

func networkConfig(cfg *config.Config) *network.Config {
	return &network.Config{
		ParentDeviceName: cfg.ParentIface,
	}
}

func appConfig(cfg *config.Config) *application.Config {
	return &application.Config{
		RootStateDir: cfg.StateRootDir,
	}
}

func appPorts(repo ports.MicroVMRepository, prov ports.MicroVMService, es ports.EventService, is ports.IDService, ns ports.NetworkService, ims ports.ImageService, fs afero.Fs) *ports.Collection {
	return &ports.Collection{
		Repo:              repo,
		Provider:          prov,
		EventService:      es,
		IdentifierService: is,
		NetworkService:    ns,
		ImageService:      ims,
		FileSystem:        fs,
	}
}

func eventSvcFromScope(ports2 *ports.Collection) ports.EventService {
	return ports2.EventService
}

func reconcileUCFromApp(app application.App) ports.ReconcileMicroVMsUseCase {
	return app
}

func queryUCFromApp(app application.App) ports.MicroVMQueryUseCases {
	return app
}

func commandUCFromApp(app application.App) ports.MicroVMCommandUseCases {
	return app
}
