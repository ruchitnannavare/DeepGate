using System;
using System.Globalization;

namespace DeepGate.Converters;

public class GroupHeaderConverter : IValueConverter
{
    public object Convert(object value, Type targetType, object parameter, CultureInfo culture)
    {
        if (targetType == typeof(Color))
        {
            return (int)value == 1 ? Colors.White : Colors.LightGray;
        }
        else if (targetType == typeof(double)) // FontSize
        {
            return (int)value == 1 ? 20.0 : 14.0;
        }
        //else if (targetType == typeof(Thickness)) // Padding
        //{
        //    return (int)value == 1 ? new Thickness(5, 5, 5, 0) : new Thickness((int)value * 15, 5, 5, 0);
        //}

        return null;
    }

    public object ConvertBack(object value, Type targetType, object parameter, CultureInfo culture)
    {
        throw new NotImplementedException();
    }
}