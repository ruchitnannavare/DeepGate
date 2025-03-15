using LiteDB;

namespace DeepGate.Models;

public class Wallpaper
{
    [BsonId]
    public int Id { get; set; }

    public string WallHavenId { get; set; }

    public string CollectionName { get; set; }

    public WallhavenResponse Response { get; set; }
}