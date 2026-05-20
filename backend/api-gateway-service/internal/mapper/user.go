package mapper

import corepb "proto/core"

func MapUserProtoToResponse(user *corepb.User) map[string]any {
	if user == nil {
		return nil
	}

	return map[string]any{
		"id":            user.Id,
		"name":          user.Name,
		"email":         user.Email,
		"created_at":    user.CreatedAt,
		"is_specialist": user.IsSpecialist,
	}
}
