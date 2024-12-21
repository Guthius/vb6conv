package export

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/guthius/vb6conv/resx"
	"github.com/guthius/vb6conv/vb6"
)

type ProjectInfo struct {
	Name      string
	Namespace string
	Output    string
}

func Export(p *ProjectInfo, f *vb6.Form) {
	control := buildControl(f.Root)
	buildMenu(control)
	resx := resx.NewResx()
	resName := filepath.Join(p.Output, control.Name+".resx")
	exportResources(resx, control)
	hasResources := resx.Count() > 0
	resx.Save(resName)
	exportForm(p, control)
	exportFormDesigner(p, control, hasResources)
	writeProgramFile(p)
	writeProjectFile(p)
}

func writeProjectFile(p *ProjectInfo) {
	fileName := filepath.Join(p.Output, p.Name+".csproj")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf(`<Project Sdk="Microsoft.NET.Sdk.WindowsDesktop">
    <PropertyGroup>
        <GenerateAssemblyInfo>False</GenerateAssemblyInfo>
        <OutputType>WinExe</OutputType>
        <UseWindowsForms>True</UseWindowsForms>
        <PlatformTarget>x86</PlatformTarget>
        <GenerateResourceUsePreserializedResources>True</GenerateResourceUsePreserializedResources>
        <TargetFramework>net48</TargetFramework>
        <LangVersion>default</LangVersion>
		 <RootNamespace>%s</RootNamespace>
    </PropertyGroup>    

	<ItemGroup>
      <PackageReference Include="System.Resources.Extensions" Version="9.0.0" />
    </ItemGroup>
</Project>`, p.Namespace))
}

func writeProgramFile(p *ProjectInfo) {
	fileName := filepath.Join(p.Output, "Program.cs")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf(`using System;
using System.Windows.Forms;

namespace %s;

internal static class Program
{
	[STAThread]
    public static void Main()
    {
        Application.Run(new frmMainMenu());
    }
}`, p.Namespace))
}

func exportResources(res resx.Resx, f *Control) {
	for k, v := range f.Resources {
		res.Add(k, v)
	}

	for _, c := range f.Children {
		exportResources(res, c)
	}
}
