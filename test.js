return (temperature >= THRESHOLD_TEMP_IN_MOUTH ? TEMP_ABOVE_WEIGHT : 0) +
           (Math.abs(diff) > THRESHOLD_TEMP_DIFFERENCE ? TEMP_DIFF_WEIGHT : 0) +
           (diff < 0 ? TEMP_NEG_WEIGHT : TEMP_BIAS);