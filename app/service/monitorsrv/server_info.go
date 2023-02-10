package monitorsrv

import (
	"context"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"go-vea/app/model/monitor/response"
	"go-vea/util"
	"time"
)

type ServerInfoService struct{}

var ServerInfoSrv = new(ServerInfoService)

func (*ServerInfoService) GetServerInfo(ctx context.Context) (*response.ServerInfo, error) {
	hostInfo, _ := host.Info()
	cpuStat, _ := cpu.Info()
	//cpuTimes, _ := cpu.Times(false)
	cpuUsed, _ := cpu.Percent(time.Second, false)

	memStat, _ := mem.VirtualMemory()

	parts, err := disk.Partitions(true)
	var disks []*response.Disk
	for _, part := range parts {
		d, _ := disk.Usage(part.Mountpoint)
		diskInfo := response.Disk{
			DirName:     part.Mountpoint,
			FsType:      part.Fstype,
			Total:       util.ConvertFileSize(d.Total),
			Used:        util.ConvertFileSize(d.Used),
			Free:        util.ConvertFileSize(d.Free),
			UsedPercent: util.DecimalPercent(d.UsedPercent),
		}
		disks = append(disks, &diskInfo)
	}

	serverInfo := response.ServerInfo{}

	localIp, ipErr := util.GetOutBoundIP()
	if ipErr != nil {
		localIp = "127.0.0.1"
	}
	sysInfo := response.Sys{
		ComputerName: hostInfo.Hostname,
		ComputerIp:   localIp,
		OsArch:       hostInfo.KernelArch,
		OsName:       hostInfo.OS,
	}
	cpuInfo := response.Cpu{
		CpuNum: cpuStat[0].Cores,
		Total:  0,
		Used:   util.Decimal(cpuUsed[0]),
		Free:   0,
	}
	memInfo := response.Mem{
		Total:       memStat.Total / (1024 * 1024 * 1024),
		Used:        memStat.Used / (1024 * 1024 * 1024),
		Free:        memStat.Free / (1024 * 1024 * 1024),
		UsedPercent: memStat.UsedPercent,
	}

	serverInfo.Sys = &sysInfo
	serverInfo.Cpu = &cpuInfo
	serverInfo.Mem = &memInfo
	serverInfo.Disk = disks
	return &serverInfo, err
}
