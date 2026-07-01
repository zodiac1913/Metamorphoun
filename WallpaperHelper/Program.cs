using System;
using System.Runtime.InteropServices;

[ComImport]
[Guid("C2CF3110-460E-4FC1-B9D0-8A1C0C9CC4BD")]
[ClassInterface(ClassInterfaceType.None)]
class DesktopWallpaperClass { }

[ComImport]
[Guid("B92B56A9-8B55-4E14-9A89-0199BBB6F93B")]
[InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IDesktopWallpaper {
    void SetWallpaper([MarshalAs(UnmanagedType.LPWStr)] string monitorID, [MarshalAs(UnmanagedType.LPWStr)] string wallpaper);
    [return: MarshalAs(UnmanagedType.LPWStr)] string GetWallpaper([MarshalAs(UnmanagedType.LPWStr)] string monitorID);
    [return: MarshalAs(UnmanagedType.LPWStr)] string GetMonitorDevicePathAt(uint monitorIndex);
    [return: MarshalAs(UnmanagedType.U4)] uint GetMonitorDevicePathCount();
    void GetMonitorRECT([MarshalAs(UnmanagedType.LPWStr)] string monitorID, out RECT displayRect);
    void SetBackgroundColor([MarshalAs(UnmanagedType.U4)] uint color);
    [return: MarshalAs(UnmanagedType.U4)] uint GetBackgroundColor();
    void SetPosition([MarshalAs(UnmanagedType.I4)] int position);
    [return: MarshalAs(UnmanagedType.I4)] int GetPosition();
    void SetSlideshow(IntPtr items);
    void GetSlideshow(out IntPtr items);
    void SetSlideshowOptions(int options, uint slideshowTick);
    void GetSlideshowOptions(out int options, out uint slideshowTick);
    void AdvanceSlideshow([MarshalAs(UnmanagedType.LPWStr)] string monitorID, [MarshalAs(UnmanagedType.I4)] int direction);
    void GetStatus(out int state);
    void Enable([MarshalAs(UnmanagedType.Bool)] bool enable);
}

[StructLayout(LayoutKind.Sequential)]
struct RECT { public int Left, Top, Right, Bottom; }

class Program {
    static int Main(string[] args) {
        if (args.Length < 1) {
            Console.Error.WriteLine("Usage: WallpaperHelper.exe <path1> [path2 ...]");
            return 1;
        }

        try {
            var wallpaper = (IDesktopWallpaper)new DesktopWallpaperClass();
            uint count = wallpaper.GetMonitorDevicePathCount();

            if (count == 0) {
                Console.Error.WriteLine("No monitors were reported by IDesktopWallpaper.");
                return 1;
            }

            if (args.Length < count) {
                Console.Error.WriteLine($"Not enough wallpaper paths for monitors. Monitors={count}, Paths={args.Length}");
                return 1;
            }

            if (args.Length > count) {
                Console.WriteLine($"Warning: received extra wallpaper paths. Monitors={count}, Paths={args.Length}. Extra paths will be ignored.");
            }

            for (uint i = 0; i < count; i++) {
                string monitorID = wallpaper.GetMonitorDevicePathAt(i);
                string path = args[i];
                wallpaper.SetWallpaper(monitorID, path);
                Console.WriteLine($"Monitor {i} [{monitorID}]: {path}");
            }
            return 0;
        } catch (Exception ex) {
            Console.Error.WriteLine($"Error: {ex.Message}");
            return 1;
        }
    }
}
