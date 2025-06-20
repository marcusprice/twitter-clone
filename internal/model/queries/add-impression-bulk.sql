UPDATE Post SET impressions = impressions + 1 WHERE id in (%s);
