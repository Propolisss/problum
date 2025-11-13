import React from 'react';
import { Tooltip } from 'react-tooltip';
import CalendarHeatmap, {
    type ReactCalendarHeatmapValue,
    type TooltipDataAttrs,
} from 'react-calendar-heatmap';
import type { ActivityData } from '../utils/activityProcessor';

import 'react-calendar-heatmap/dist/styles.css';

type Props = {
    data: ActivityData;
};

type HeatmapValue = ReactCalendarHeatmapValue<Date>;

export default function ActivityHeatmap({ data }: Props) {
    if (!data) return null;

    const values = data.heatmapValues as HeatmapValue[];

    return (
        <div>
            <div className="flex items-center gap-6 mb-4">
                <div>
                    <div className="text-sm text-gray-500">Активных дней</div>
                    <div className="text-2xl font-bold">{data.totalActiveDays}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-500">Макс. стрик</div>
                    <div className="text-2xl font-bold">{data.maxStreak}</div>
                </div>
            </div>

            <CalendarHeatmap<Date>
                startDate={data.startDate}
                endDate={data.endDate}
                values={values}
                classForValue={(value?: HeatmapValue) => {
                    if (!value) return 'color-empty';
                    const count = Math.min(4, value.count ?? 0);
                    return `color-scale-${count}`;
                }}
                tooltipDataAttrs={(value?: HeatmapValue): TooltipDataAttrs => {
                    const content =
                        value && value.date
                            ? `${value.count ? `${value.count} решений` : 'Нет решений'} • ${value.date.toLocaleDateString()}`
                            : '';

                    const attrs: Record<string, string> = {
                        ['data-tooltip-id']: 'heatmap-tooltip',
                        ['data-tooltip-content']: content,
                    };
                    return attrs as unknown as TooltipDataAttrs;
                }}
                showMonthLabels={true}
            />

            <Tooltip id="heatmap-tooltip" />
        </div>
    );
}
