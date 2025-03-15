using System;
using LiteDB;

namespace DeepGate.Models
{
    public class BackgroundSetting
    {
        [BsonId]
        public int Id { get; set; }
        
        public int SelectedWallpaperIndex { get; set; }
        
        public Wallpaper SelectedWallpaper { get; set; }

        public int Opacity { get; set; }
    }
} 