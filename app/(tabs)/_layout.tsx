import { Tabs } from 'expo-router';
import React from 'react';
import { GarageIcon, CarIcon, SettingsIcon, MaintenanceIcon, ReminderIcon } from '@/components/ui/IconSymbol';

import { Colors } from '@/constants/theme';
import { useColorScheme } from '@/hooks/use-color-scheme';

export default function TabLayout() {
  const colorScheme = useColorScheme();

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: Colors[colorScheme ?? 'light'].tint,
        tabBarStyle: {
          height: 80,
        },
        tabBarItemStyle: {
          justifyContent: 'center',
          alignItems: 'center'
        },
        headerShown: false,
      }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
         tabBarIcon: ({ color, size }) => <GarageIcon size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="Car"
        options={{
          title: 'Car',
        tabBarIcon: ({ color, size }) => <CarIcon size={size} color={color} />,

        }}
      />
      <Tabs.Screen
        name="Maintenance"
        options={{
          title: 'Maintenance',
        tabBarIcon: ({ color, size }) => <MaintenanceIcon size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="Reminder"
        options={{
          title: 'Reminders',
        tabBarIcon: ({ color, size }) => <ReminderIcon size={size} color={color} />,
        }}
      />
      <Tabs.Screen
        name="Settings"
        options={{
          title: 'Settings',
        tabBarIcon: ({ color, size }) => <SettingsIcon size={size} color={color} />,
        }}
      />
    </Tabs>
  );
}
